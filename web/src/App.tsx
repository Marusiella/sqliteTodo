import "@picocss/pico";
import { useEffect, useState } from "react";
import "./App.css";
function App() {
  const [todos, setTodos] = useState<{content: string, id: number}[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [text, setText] = useState<string>("");
  const fetchTodos = async () => {
    setLoading(true);
    try {
      const res = await fetch("/todos");
      if (res.ok) {
        const json = await res.json();
        setTodos(json);
      }
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    fetchTodos();
  }, []);
  useEffect(() => {
    if (text === "") {
      fetchTodos();
    }
    fetch("/todos/search", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(text),
    })
      .then((res) => {
        if (res.ok) {
          return res.json();
        }
      })
      .then((json) => {
        if (json === null) return;
        setTodos(json.map((todo: { content: string, id: number }) => todo));
      });
  }, [text]);
  const handleSubbmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const res = await fetch("/todos", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(text),
    });

    if (res.ok) {
      const json = await res.json();
      setTodos([...todos, json]);
      setText("");
    }
  };
  const handleDelete = async (e: string) => {
    const res = await fetch("/todos/" + e, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (res.ok) {
      fetchTodos();
    }
  };
  return (
    <div
      className="container"
      style={{ paddingTop: "40px", paddingBottom: "40px" }}
    >
      <form
        onSubmit={handleSubbmit}
        className="container"
        style={{ display: "flex" }}
      >
        <input
          type="text"
          onChange={(e) => setText(e.target.value)}
          value={text}
        />
        <button type="submit" style={{ width: "20%", marginLeft: "20px" }}>
          Add
        </button>
      </form>
      {loading && <article aria-busy="true"></article>}
      <main>
        {todos.map((todo, i) => (
          <article
            className="card flex-end"
            style={{ animationDelay: `${i <= 5 ? i * 100 : 5 * 100}ms` }}
            key={todo.id}
          >
            {" "}
            <div
              dangerouslySetInnerHTML={{
                __html: todo.content,
              }}
            />
            <button
              className="secondary outline"
              style={{ marginLeft: "20px", maxWidth: "130px", color: "light" }}
              onClick={() => handleDelete(todo.id.toString())}
            >
              Delete
            </button>
          </article>
        ))}
      </main>
    </div>
  );
}

export default App;
