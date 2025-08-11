import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/books/$booksId')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/books/$booksId"!</div>
}
