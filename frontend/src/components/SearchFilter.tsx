import React, { useState, useEffect } from 'react';

interface SearchFilterProps {
  onSearch: (query: string) => void;
  initialValue?: string;
}

export const SearchFilter: React.FC<SearchFilterProps> = ({ onSearch, initialValue = '' }) => {
  const [query, setQuery] = useState(initialValue);

  useEffect(() => {
    const debounceTimer = setTimeout(() => {
      onSearch(query);
    }, 300); // Debounce search by 300ms

    return () => clearTimeout(debounceTimer);
  }, [query, onSearch]);

  return (
    <div className="mb-4">
      <div className="relative">
        <input
          type="text"
          placeholder="Search URLs or page titles..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="w-full px-4 py-2 pl-10 pr-4 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <svg className="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
      </div>
    </div>
  );
};