package energy

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nhuongmh/cvfs/timemn/bootstrap"
	"github.com/nhuongmh/cvfs/timemn/internal/logger"
	"github.com/nhuongmh/cvfs/timemn/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/api/sheets/v4"
)

const (
	ENERGY_DATE_COL        = 0
	ENERGY_SLEEP_START_COL = 1
	ENERGY_SLEEP_END_COL   = 2
	ENERGY_EXECISE_COL     = 3
	ENERGY_NUTS_COL        = 4
	ENERGY_SXS_COL         = 5
	ENERGY_FEELING_ADJ_COL = 6
	ENERGY_SLEEP_SCORE_COL = 7
	ENERGY_ETF_COL         = 8
	ENERGY_ETD_COL         = 9
	ENERGY_ETD_PERC_COL    = 10
	ENERGY_NOTE_COL        = 11
	ENERGY_SYSNOTE_COL     = 12
)

type energyMngService struct {
	GgSheetSrv *sheets.Service
	env        *bootstrap.Env
}

func NewEnergyMngService(env *bootstrap.Env) model.EnergyMngService {
	logger.Log.Info().Msg("Initializing energy management service")
	ggSheetSrv, err := InitNewGoogleSheetService(env.GoogleKeyBase64)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize google sheet service")
	}
	return &energyMngService{
		GgSheetSrv: ggSheetSrv,
		env:        env,
	}
}

func (es *energyMngService) EvaluateAllFromSheet(ctx context.Context, forceOverwrite bool) error {
	dailyInputList, err := es.fetchAllDataFromSheet(es.env.GoogleSpreadSheetId, es.env.GoogleEnergySheetName)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch data from sheet")
	}
	logger.Log.Info().Msgf("Fetched %v rows from sheet", len(*dailyInputList))

	for i := range *dailyInputList {
		dailyInput := (*dailyInputList)[i]
		logger.Log.Info().Msgf("Evaluating row %v, date=%v", dailyInput.Row, dailyInput.Date)
		sleepScore, etf, err := es.evaluateEtf(&dailyInput)

		if err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to evaluate etf for row %v, date=%v", dailyInput.Row, dailyInput.Date)
			continue
		}

		if etf < 0 {
			logger.Log.Warn().Msgf("Invalid etf for row %v, date=%v", dailyInput.Row, dailyInput.Date)
			etf = 0
		}

		if err := es.handleEtfMismatch(&dailyInput, etf, forceOverwrite); err != nil {
			return err
		}

		if err := es.writeValueToCell(dailyInput.Row, ENERGY_ETF_COL, fmt.Sprintf("%.2f", etf)); err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to write etf to row %v", dailyInput.Row)
			continue
		}

		if err := es.handleSleepScoreMismatch(&dailyInput, sleepScore, forceOverwrite); err != nil {
			return err
		}

		if err := es.writeValueToCell(dailyInput.Row, ENERGY_SLEEP_SCORE_COL, fmt.Sprintf("%.2f", sleepScore)); err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to write sleep score to row %v", dailyInput.Row)
			continue
		}

		etd := es.calculateEtd(i, dailyInputList, etf)
		if err := es.writeValueToCell(dailyInput.Row, ENERGY_ETD_COL, fmt.Sprintf("%.2f", etd)); err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to write etd to row %v", dailyInput.Row)
			continue
		}
		dailyInput.Etd = etd
		if err := es.writeValueToCell(dailyInput.Row, ENERGY_ETD_PERC_COL, fmt.Sprintf("%.0f%%", etd*100)); err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to write etd percent to row %v", dailyInput.Row)
			continue
		}
	}

	return nil
}

func (es *energyMngService) handleEtfMismatch(dailyInput *model.DailyPhysicalInput, etf float64, forceOverwrite bool) error {
	if dailyInput.Etf != 0 && dailyInput.Etf != etf {
		logger.Log.Warn().Msgf("Mismatch Etf in row %v, date=%v, existed=%v, new calculation=%v", dailyInput.Row, dailyInput.Date, dailyInput.Etf, etf)
		if !forceOverwrite {
			dailyInput.Note = "System err: Etf mismatch"
			es.writeValueToCell(dailyInput.Row, ENERGY_SYSNOTE_COL, dailyInput.Note)
			return model.ErrSystemCalculationError
		}
	}
	return nil
}

func (es *energyMngService) handleSleepScoreMismatch(dailyInput *model.DailyPhysicalInput, sleepScore float64, forceOverwrite bool) error {
	if dailyInput.SleepScore != 0 && dailyInput.SleepScore != sleepScore {
		logger.Log.Warn().Msgf("Mismatch SleepScore in row %v, date=%v, existed=%v, new calculation=%v", dailyInput.Row, dailyInput.Date, dailyInput.SleepScore, sleepScore)
		if !forceOverwrite {
			dailyInput.Note = "System err: SleepScore mismatch"
			es.writeValueToCell(dailyInput.Row, ENERGY_SYSNOTE_COL, dailyInput.Note)
			return model.ErrSystemCalculationError
		}
	}
	return nil
}

func (es *energyMngService) calculateEtd(index int, dailyInputList *[]model.DailyPhysicalInput, etf float64) float64 {
	if index > 0 {
		return es.evaluateEtd((*dailyInputList)[index-1].Etd, etf)
	}
	return etf
}

func (es *energyMngService) evaluateEtf(input *model.DailyPhysicalInput) (sleepScore float64, etf float64, errx error) {
	if model.SLEEPING_WEIGHT+model.EXERCISE_WEIGHT+model.HS_WEIGHT+model.NUTS_WEIGHT+model.FEELING_WEIGHT != 1 {
		logger.Log.Warn().Msg("Weights do not sum up to 1")
		return 0, 0, model.ErrInvalidInput
	}

	exerciseScore := 0.0
	if len(input.Exercise) > 0 {
		totalDuration := 0.0
		for _, duration := range input.Exercise {
			totalDuration += duration
		}
		exerciseScore = totalDuration / 0.5
		if exerciseScore > 1 {
			exerciseScore = 1
		}
	}

	//get sleep duration
	sleepDuration := input.Sleep.EndSleepingTime.Sub(input.Sleep.StartSleepingTime).Hours()
	if sleepDuration < 0 {
		// Handle case where sleep end time is on the next day
		sleepDuration += 24
	}
	if sleepDuration < 0 {
		logger.Log.Warn().Msg("Invalid sleep start/end time")
		return 0, 0, model.ErrInvalidInput
	}
	startHour := input.Sleep.StartSleepingTime.Hour()
	startMin := input.Sleep.StartSleepingTime.Minute()
	if startHour > 18 {
		startHour -= 24
	}
	startHour += startMin / 60

	sleepScore = sleepDuration/7 - 0.15*(float64(startHour)+1)
	if sleepScore > 1 {
		sleepScore = 1
	}
	logger.Log.Info().Msgf("Sleep duration: %v, startHour=%v, Score=%v", sleepDuration, startHour, sleepScore)

	nfScore := 1.0
	if input.Nuts > 0 {
		nfScore = 0.0
	}

	nsScore := 1.0
	if input.Sxs > 0 {
		nsScore = 0.0
	}

	if input.Feeling < 0 || input.Feeling > 1 {
		logger.Log.Warn().Msg("Invalid feeling adjustment")
		return 0, 0, model.ErrInvalidInput
	}

	etf = model.SLEEPING_WEIGHT*sleepScore +
		model.EXERCISE_WEIGHT*exerciseScore +
		model.NUTS_WEIGHT*nfScore +
		model.HS_WEIGHT*nsScore +
		model.FEELING_WEIGHT*input.Feeling

	logger.Log.Info().Msgf("SleepScore=%v, ExerciseScore=%v, NfScore=%v, NsScore=%v, Feeling=%v", sleepScore, exerciseScore, nfScore, nsScore, input.Feeling)
	return sleepScore, etf, nil
}

func (es *energyMngService) evaluateEtd(yesterdayEtd float64, todayEtf float64) float64 {
	if yesterdayEtd == 0 {
		return todayEtf
	}
	return (2*todayEtf + yesterdayEtd) / 3
}

func (es *energyMngService) writeValueToCell(row, col int, value string) error {
	valueRange := fmt.Sprintf("%s!%s%d", es.env.GoogleEnergySheetName, string(rune('A'+col)), row)
	valueInputOption := "RAW"
	valueRangeBody := &sheets.ValueRange{
		Values: [][]interface{}{{value}},
	}
	_, err := es.GgSheetSrv.Spreadsheets.Values.Update(es.env.GoogleSpreadSheetId, valueRange, valueRangeBody).ValueInputOption(valueInputOption).Do()
	if err != nil {
		return errors.Wrapf(err, "Failed to write value to cell %v", valueRange)
	}
	return nil
}

func (es *energyMngService) fetchAllDataFromSheet(spreadsheetId, sheetName string) (*[]model.DailyPhysicalInput, error) {
	// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>

	logger.Log.Info().Msgf("Fetching data from google sheet %v (%v)", spreadsheetId, sheetName)
	readRange := fmt.Sprintf("%s!A2:L", sheetName)
	resp, err := es.GgSheetSrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve data from sheet")
	}

	if len(resp.Values) == 0 {
		logger.Log.Warn().Msg("No data found.")
		return nil, model.ErrNoData
	}

	logger.Log.Info().Msgf("Fetched %v rows from sheet", len(resp.Values))

	dailyEList := make([]model.DailyPhysicalInput, 0, 20)
	for idx, row := range resp.Values {
		if len(row) < 7 {
			logger.Log.Warn().Msgf("Invalid row %v, not enough data. STOP parsing", idx+2)
			break
		}

		realRow := func(rowIdx int) int {
			return rowIdx + 2
		}

		getOptional := func(col int) string {
			if len(row) > col {
				return row[col].(string)
			}
			return ""
		}

		date := row[ENERGY_DATE_COL].(string)              //1
		sleepStart := row[ENERGY_SLEEP_START_COL].(string) //2
		sleepEnd := row[ENERGY_SLEEP_END_COL].(string)     //3
		exercise := row[ENERGY_EXECISE_COL].(string)       //4
		nuts := row[ENERGY_NUTS_COL].(string)              //5
		sxs := row[ENERGY_SXS_COL].(string)                //6
		feelingAdj := row[ENERGY_FEELING_ADJ_COL].(string) //7

		sleepScore := getOptional(ENERGY_SLEEP_SCORE_COL)
		etf := getOptional(ENERGY_ETF_COL)
		etd := getOptional(ENERGY_ETD_COL)
		note := getOptional(ENERGY_NOTE_COL)

		parsedDate, err := parseDate(date)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to parse date %v in row %v", date, realRow(idx))
			continue
		}
		parsedSleep, err := parseSleepTime(sleepStart, sleepEnd)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to parse sleep time %v-%v in row %v", sleepStart, sleepEnd, realRow(idx))
			continue
		}
		parsedExercise, err := parseExercise(exercise)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to parse exercise %v in row %v", exercise, realRow(idx))
			continue
		}
		parsedNuts, err := strconv.Atoi(nuts)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to parse nuts %v in row %v", nuts, realRow(idx))
			continue
		}

		parsedSxs, err := strconv.Atoi(sxs)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to parse sxs %v in row %v", sxs, realRow(idx))
			continue
		}

		parsedFeelingAdj, err := strconv.ParseFloat(feelingAdj, 64)
		if err != nil {
			logger.Log.Warn().Err(err).Msgf("Failed to parse feeling adj %v in row %v", feelingAdj, realRow(idx))
			parsedFeelingAdj = 0.5
		}

		parsedSleepScore, err := strconv.ParseFloat(sleepScore, 64)
		if err != nil {
			parsedSleepScore = 0
		}

		parsedEtf, err := strconv.ParseFloat(etf, 64)
		if err != nil {
			parsedEtf = 0
		}

		parsedEtd, err := strconv.ParseFloat(etd, 64)
		if err != nil {
			parsedEtd = 0
		}

		dailyEList = append(dailyEList, model.DailyPhysicalInput{
			Row:        idx + 2,
			Date:       parsedDate,
			Sleep:      parsedSleep,
			Nuts:       parsedNuts,
			Sxs:        parsedSxs,
			Exercise:   *parsedExercise,
			Feeling:    parsedFeelingAdj,
			SleepScore: parsedSleepScore,
			Etf:        parsedEtf,
			Etd:        parsedEtd,
			Note:       note,
		})

	}

	return &dailyEList, nil
}
