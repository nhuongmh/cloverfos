"use client";
import React, { useState } from "react";
import { FiSearch, FiPlus } from "react-icons/fi";
import { IoMdClose } from "react-icons/io";

// Mock Data
const mockCategories = [
  {
    id: 1,
    title: "Personal Growth",
    description: "Discover resources to enhance your personal development journey",
    image: "images.unsplash.com/photo-1544367567-0f2fcb009e0b"
  },
  {
    id: 2,
    title: "Career Success",
    description: "Strategies and inspiration for professional excellence",
    image: "images.unsplash.com/photo-1522202176988-66273c2fd55f"
  },
  {
    id: 3,
    title: "Health & Wellness",
    description: "Motivational content for maintaining a healthy lifestyle",
    image: "images.unsplash.com/photo-1517836357463-d25dfeac3438"
  },
  {
    id: 4,
    title: "Mindfulness",
    description: "Resources for mental clarity and emotional balance",
    image: "images.unsplash.com/photo-1506126613408-eca07ce68773"
  }
];

const MotivationPacks = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);
  const [newCategory, setNewCategory] = useState({
    title: "",
    description: "",
    image: ""
  });

  const handleSearch = (e) => {
    setSearchQuery(e.target.value);
  };

  const filteredCategories = mockCategories.filter(
    (category) =>
      category.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      category.description.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleSubmit = (e) => {
    e.preventDefault();
    // Handle form submission logic here
    setShowAddForm(false);
    setNewCategory({ title: "", description: "", image: "" });
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6">
      <div className="max-w-7xl mx-auto">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-4">Motivation Packs</h1>
          <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
            <div className="relative w-full md:w-96">
              <FiSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="Search categories..."
                className="w-full bg-gray-800 rounded-lg py-2 pl-10 pr-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={searchQuery}
                onChange={handleSearch}
              />
            </div>
            <button
              onClick={() => setShowAddForm(true)}
              className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded-lg transition-colors"
            >
              <FiPlus /> Add New Category
            </button>
          </div>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCategories.map((category) => (
            <div
              key={category.id}
              className="bg-gray-800 rounded-lg overflow-hidden hover:shadow-lg transition-shadow"
            >
              <img
                src={`https://${category.image}`}
                alt={category.title}
                className="w-full h-48 object-cover"
                onError={(e) => {
                  e.target.onerror = null;
                  e.target.src = "https://images.unsplash.com/photo-1516383740770-fbcc5ccbece0";
                }}
              />
              <div className="p-4">
                <h3 className="text-xl font-semibold mb-2">{category.title}</h3>
                <p className="text-gray-400">{category.description}</p>
              </div>
            </div>
          ))}
        </div>

        {showAddForm && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div className="bg-gray-800 rounded-lg p-6 w-full max-w-md">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-2xl font-bold">Add New Category</h2>
                <button
                  onClick={() => setShowAddForm(false)}
                  className="text-gray-400 hover:text-white"
                >
                  <IoMdClose size={24} />
                </button>
              </div>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block mb-1">Title</label>
                  <input
                    type="text"
                    required
                    className="w-full bg-gray-700 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={newCategory.title}
                    onChange={(e) =>
                      setNewCategory({ ...newCategory, title: e.target.value })
                    }
                  />
                </div>
                <div>
                  <label className="block mb-1">Description</label>
                  <textarea
                    required
                    className="w-full bg-gray-700 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows="3"
                    value={newCategory.description}
                    onChange={(e) =>
                      setNewCategory({ ...newCategory, description: e.target.value })
                    }
                  ></textarea>
                </div>
                <div>
                  <label className="block mb-1">Image URL</label>
                  <input
                    type="url"
                    required
                    className="w-full bg-gray-700 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={newCategory.image}
                    onChange={(e) =>
                      setNewCategory({ ...newCategory, image: e.target.value })
                    }
                  />
                </div>
                <button
                  type="submit"
                  className="w-full bg-blue-600 hover:bg-blue-700 py-2 rounded-lg transition-colors"
                >
                  Add Category
                </button>
              </form>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default MotivationPacks;