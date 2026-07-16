import { useEffect, useState } from 'react'
import './App.css'

const API_BASE = 'http://localhost:8080/api'

function App() {
  const [todos, setTodos] = useState([])
  const [title, setTitle] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  async function loadTodos() {
    try {
      const res = await fetch(`${API_BASE}/todos`)
      if (!res.ok) throw new Error('Gagal mengambil data')
      setTodos(await res.json())
      setError(null)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTodos()
  }, [])

  async function addTodo(e) {
    e.preventDefault()
    if (!title.trim()) return
    await fetch(`${API_BASE}/todos`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title }),
    })
    setTitle('')
    loadTodos()
  }

  async function toggleTodo(id) {
    await fetch(`${API_BASE}/todos/${id}`, { method: 'PATCH' })
    loadTodos()
  }

  async function deleteTodo(id) {
    await fetch(`${API_BASE}/todos/${id}`, { method: 'DELETE' })
    loadTodos()
  }

  return (
    <div className="app">
      <h1>Todo App — React + Go</h1>
      <p className="status">
        {loading ? 'Memuat...' : error ? `Error: ${error}` : `${todos.length} task`}
      </p>

      <form className="todo-form" onSubmit={addTodo}>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Tambah task baru..."
        />
        <button type="submit">Tambah</button>
      </form>

      <ul className="todo-list">
        {todos.map((t) => (
          <li key={t.id} className={`todo-item ${t.done ? 'done' : ''}`}>
            <input type="checkbox" checked={t.done} onChange={() => toggleTodo(t.id)} />
            <span>{t.title}</span>
            <button onClick={() => deleteTodo(t.id)}>Hapus</button>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default App
