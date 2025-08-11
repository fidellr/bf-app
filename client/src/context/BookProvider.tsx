import { createContext, useContext } from "react";

import type { Book } from "@/types/books";

type BookContextType = {
  books: Book[];
  loading: boolean;
  error: string | null;
  createBook: (book: Omit<Book, "id">) => Promise<void>;
  updateBook: (id: string, book: Partial<Book>) => Promise<void>;
  deleteBook: (id: string) => Promise<void>;
};

export const BookContext = createContext<BookContextType | undefined>(
  undefined
);

export const useBookContext = () => {
  const context = useContext(BookContext);
  if (context === undefined) {
    throw new Error("useBookContext must be used within a BookProvider");
  }
  return context;
};
