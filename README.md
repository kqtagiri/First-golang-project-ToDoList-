# First-golang-project-ToDoList-
In version 2.1 i change type of my list. In the last version 2.0 list have type  slice. Now he have type map with name in index. With this update i can easily search task that i need or check, have been another task before with this title.
I also create my router and split my handler to many handlers. Router help me with the distribution of handlers.
And i add Mutex and RWMutex in my handlers to escape race condition.

In version 3.0 i update my todo list with adding database. I use database postgreSQL for saving my tasks. And now, when the programm end, all tasks, which user created, save in database and user can use it, when he start programm then.