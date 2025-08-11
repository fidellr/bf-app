import type { Book } from "@/types/books";

export type BookCreate = Omit<Book, "id">;
export type BookUpdate = Partial<BookCreate>;

export interface BookContextType {
  books: Book[];
  loading: boolean;
  error: string | null;
  fetchBooks: () => Promise<void>;
  getBook: (id: string) => Book | undefined;
  createBook: (book: BookCreate) => Promise<void>;
  updateBook: (id: string, updates: BookUpdate) => Promise<void>;
  deleteBook: (id: string) => Promise<void>;
}
