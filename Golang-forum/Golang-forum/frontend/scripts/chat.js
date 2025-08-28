let socket;
const wsStatus = document.getElementById("ws-status");
const isAdmin = localStorage.getItem("isAdmin") === "true";

function connectWebSocket() {
  const token = localStorage.getItem("token");
  const wsUrl = token
    ? `ws://localhost:8081/ws?token=${token}`
    : "ws://localhost:8081/ws";

  socket = new WebSocket(wsUrl);

  socket.onopen = () => {
    wsStatus.textContent = "Подключено";
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

    if (Array.isArray(data)) {
      document.getElementById("messages").innerHTML = "";
      data.reverse().forEach(addMessage);
    } else {
      addMessage(data);
    }

    document.getElementById("messages").scrollTop =
      document.getElementById("messages").scrollHeight;
  };

  socket.onerror = () => {
    wsStatus.textContent = "Ошибка соединения";
  };

  socket.onclose = () => {
    wsStatus.textContent = "Отключено";
  };
}

function addMessage(msg) {
  const createdAt = msg.created_at || msg.createdAt;
  const messageDiv = document.createElement("div");
  messageDiv.className = "message";
  messageDiv.innerHTML = `
    <span class="author">${msg.username}</span>
    <span class="text">${msg.text}</span>
    <span class="time">${
      createdAt ? new Date(createdAt).toLocaleTimeString() : ""
    }</span>
  `;

  if (isAdmin) {
    const deleteBtn = document.createElement("button");
    deleteBtn.textContent = "Удалить";
    deleteBtn.className = "delete-btn";

    deleteBtn.onclick = async () => {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          alert("Требуется авторизация");
          return;
        }
        const res = await fetch(
          `http://localhost:8080/api/admin/delete-message`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: "Bearer " + token,
            },
            body: JSON.stringify({
              message_id: msg.id,
              is_admin: isAdmin,
            }),
          }
        );
        if (!res.ok) {
          const text = await res.text();
          alert("Ошибка удаления сообщения: " + text);
        } else {
          messageDiv.remove(); // Успешно → удаляем сообщение из DOM
        }
      } catch (err) {
        console.error("Ошибка при удалении сообщения:", err);
        alert("Не удалось отправить запрос на удаление сообщения.");
      }
    };
    messageDiv.appendChild(deleteBtn);
  }
  document.getElementById("messages").appendChild(messageDiv);
}

if (isAdmin) {
  document.getElementById("admin-panel").style.display = "block";
  loadUsers();
}

async function loadUsers() {
  const res = await fetch("http://localhost:8080/api/admin/users", {
    headers: {
      Authorization: "Bearer " + localStorage.getItem("token"),
    },
  });
  const users = await res.json();
  const tbody = document.querySelector("#user-table tbody");
  tbody.innerHTML = "";
  users.forEach((user) => {
    const row = document.createElement("tr");
    row.innerHTML = `
      <td>${user.username}</td>
      <td>${user.email}</td>
      <td>${user.banned ? "Забанен" : "Активен"}</td>
      <td>
        <button onclick="toggleBan('${user.email}', ${user.banned})">
          ${user.banned ? "Разблокировать" : "Забанить"}
        </button>
      </td>
    `;
    tbody.appendChild(row);
  });
}

async function toggleBan(email, isBanned) {
  const endpoint = isBanned ? "unban" : "ban";
  const res = await fetch(`http://localhost:8080/api/admin/${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + localStorage.getItem("token"),
    },
    body: JSON.stringify({ email }),
  });
  if (!res.ok) {
    const text = await res.text();
    alert("Ошибка изменения статуса пользователя: " + text);
  }
  loadUsers();
}

document.addEventListener("DOMContentLoaded", () => {
  connectWebSocket();

  const token = localStorage.getItem("token");
  const email = localStorage.getItem("email");
  const username = localStorage.getItem("username");
  const profileDiv = document.getElementById("profile");

  if (token && email && username) {
    profileDiv.textContent = `Вы вошли как: ${username} (${email})`;
  } else {
    profileDiv.textContent = `Вы зашли как гость (только чтение)`;

    // Отключаем ввод сообщения и кнопку отправки
    document.getElementById("message-input").disabled = true;
    document.getElementById("send-btn").disabled = true;
  }

  if (!isAdmin) {
    const adminPanel = document.getElementById("admin-panel");
    if (adminPanel) {
      adminPanel.remove();
    }
  }
});

// Отправка сообщений
document.getElementById("send-btn").addEventListener("click", () => {
  const token = localStorage.getItem("token");
  if (!token) {
    alert("Вы должны войти в аккаунт, чтобы отправлять сообщения.");
    return;
  }

  const input = document.getElementById("message-input");
  const text = input.value.trim();
  if (!text) return;

  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ text }));
    input.value = "";
  } else {
    console.error("WebSocket is not connected");
  }
});

document.getElementById("message-input").addEventListener("keydown", (e) => {
  if (e.key === "Enter") {
    document.getElementById("send-btn").click();
  }
});

document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("token");
  localStorage.removeItem("email");
  localStorage.removeItem("username");
  localStorage.removeItem("isAdmin");
  window.location.href = "login.html";
});
