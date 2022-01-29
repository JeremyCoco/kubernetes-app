package mysql

const (
	insertUser = "INSERT INTO user(username, hashed_password) VALUES(?, ?);"

	insertTodo = "INSERT INTO todo(title, created_at, expires) VALUES(?, ?, ?);"

	insertUserTodo = "INSERT INTO user_todo(user_id, todo_id) VALUES(?, ?);"

	selectUserForAuth = `SELECT id, hashed_password FROM user
    WHERE username = ?;`

	selectUserById = `SELECT id, username FROM user WHERE id = ?;`

	selectAllTodos = `SELECT todo.id, todo.title, todo.done, todo.created_at, todo.expires
    FROM user_todo
    INNER JOIN user
    ON user_todo.user_id = user.id
    INNER JOIN todo
    ON user_todo.todo_id = todo.id
    WHERE user.id = ?
    ORDER BY created_at DESC;
    `
	selectUserTodo = "SELECT todo_id FROM user_todo WHERE user_id = ? AND todo_id = ?;"

	updateTodoToDone = "UPDATE todo SET done = TRUE WHERE id = ?;"

	deleteTodo = "DELETE FROM todo WHERE id = ?;"

	deleteUserTodo = "DELETE FROM user_todo WHERE user_id = ? AND todo_id = ?;"
)
