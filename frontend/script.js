const taskForm = document.getElementById("taskForm");
const taskList = document.getElementById("taskList");

// CREATE TASK
taskForm.addEventListener("submit", async function (event) {

    event.preventDefault();

    const title = document.getElementById("title").value;
    const description = document.getElementById("description").value;

    const task = {
        title: title,
        description: description
    };

    try {

        const response = await fetch("/tasks", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(task)
        });

        if (!response.ok) {
            throw new Error("Failed to create task");
        }

        taskForm.reset();

        loadTasks();

    } catch (error) {

        console.error(error);
        alert("Could not create task");

    }
});


// READ TASKS
async function loadTasks() {

    try {

        const response = await fetch("/tasks");

        if (!response.ok) {
            throw new Error("Failed to fetch tasks");
        }

        const tasks = await response.json();

        taskList.innerHTML = "";

        tasks.forEach(function (task) {

            const taskElement = document.createElement("div");

            taskElement.className = "task";

            taskElement.innerHTML = `
                <h3>${task.title}</h3>
                <p>${task.description}</p>

                <button onclick="editTask(${task.id}, '${task.title}', '${task.description}')">
                    Edit
                </button>

                <button onclick="deleteTask(${task.id})">
                    Delete
                </button>
            `;

            taskList.appendChild(taskElement);
        });

    } catch (error) {

        console.error(error);
        taskList.innerHTML = "<p>Could not load tasks.</p>";

    }
}


// UPDATE TASK
async function editTask(id, oldTitle, oldDescription) {

    const newTitle = prompt("Enter new title:", oldTitle);

    if (newTitle === null) {
        return;
    }

    const newDescription = prompt(
        "Enter new description:",
        oldDescription
    );

    if (newDescription === null) {
        return;
    }

    const updatedTask = {
        title: newTitle,
        description: newDescription
    };

    try {

        const response = await fetch(`/tasks/${id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(updatedTask)
        });

        if (!response.ok) {
            throw new Error("Failed to update task");
        }

        loadTasks();

    } catch (error) {

        console.error(error);
        alert("Could not update task");

    }
}


// DELETE TASK
async function deleteTask(id) {

    const confirmDelete = confirm(
        "Are you sure you want to delete this task?"
    );

    if (!confirmDelete) {
        return;
    }

    try {

        const response = await fetch(`/tasks/${id}`, {
            method: "DELETE"
        });

        if (!response.ok) {
            throw new Error("Failed to delete task");
        }

        loadTasks();

    } catch (error) {

        console.error(error);
        alert("Could not delete task");

    }
}


// LOAD TASKS WHEN PAGE OPENS
loadTasks();