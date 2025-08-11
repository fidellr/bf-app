import React, { createContext, useContext, useState, useEffect } from "react";
import type { BookContextType, BookCreate } from "@/context/BookTypes";
import { getBooks, createBook, updateBook, deleteBook } from "@/api/books";
import type { Book } from "@/types/books";

const BookContext = createContext<BookContextType | undefined>(undefined);

export const BookProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchBooks = async () => {
    setLoading(true);
    try {
      const data = await getBooks();
      setBooks(data);
      setError(null);
    } catch (err) {
      setError("Failed to fetch books");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const getBook = (id: number) => books.find((book) => book.id === id);

  const handleCreate = async (book: BookCreate) => {
    try {
      const newBook = await createBook(book);
      setBooks((prev) => [...prev, newBook]);
    } catch (err) {
      setError("Failed to create book");
      throw err;
    }
  };

  const handleUpdate = async (id: string, updates: BookUpdate) => {
    try {
      const updatedBook = await updateBook(id, updates);
      setBooks((prev) =>
        prev.map((book) =>
          book.id === id ? { ...book, ...updatedBook } : book
        )
      );
    } catch (err) {
      setError("Failed to update book");
      throw err;
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteBook(id);
      setBooks((prev) => prev.filter((book) => book.id !== id));
    } catch (err) {
      setError("Failed to delete book");
      throw err;
    }
  };

  useEffect(() => {
    fetchBooks();
  }, []);

  const value = {
    books,
    loading,
    error,
    fetchBooks,
    getBook,
    createBook: handleCreate,
    updateBook: handleUpdate,
    deleteBook: handleDelete,
  };

  return <BookContext.Provider value={value}>{children}</BookContext.Provider>;
};

export const useBooks = () => {
  const context = useContext(BookContext);
  if (context === undefined) {
    throw new Error("useBooks must be used within a BookProvider");
  }
  return context;
};
