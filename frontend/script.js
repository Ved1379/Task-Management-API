const taskForm = document.getElementById("taskForm");
const taskList = document.getElementById("taskList");

taskForm.addEventListener("submit", async function (event) {

    event.preventDefault();

    const title = document.getElementById("title").value;
    const description = document.getElementById("description").value;

    const task = {
        title: title,
        description: description
    };

    const response = await fetch("/tasks", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(task)
    });

    const createdTask = await response.json();

    console.log("Created task:", createdTask);

    taskForm.reset();

    loadTasks();
});


async function loadTasks() {

    const response = await fetch("/tasks");

    const tasks = await response.json();

    taskList.innerHTML = "";

    tasks.forEach(function (task) {

        const taskElement = document.createElement("div");

        taskElement.className = "task";

        taskElement.innerHTML = `
            <h3>${task.title}</h3>
            <p>${task.description}</p>
        `;

        taskList.appendChild(taskElement);
    });
}


loadTasks();