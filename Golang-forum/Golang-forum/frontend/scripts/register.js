const API_URL = "http://localhost:8080/api";
const emailInput = document.getElementById("email");
const usernameInput = document.getElementById("username");
const passwordInput = document.getElementById("password");
const isAdminInput = document.getElementById("is-admin");
const registerBtn = document.getElementById("register-btn");
const errorMessage = document.getElementById("error-message");

function validateEmail(email) {
  return /\S+@\S+\.\S+/.test(email);
}

registerBtn.addEventListener("click", async () => {
  const email = emailInput.value.trim();
  const username = usernameInput.value.trim();
  const password = passwordInput.value;
  const isAdmin = isAdminInput.checked;
  if (!validateEmail(email) || !username || !password) {
    errorMessage.textContent = "Заполните все поля корректно";
    return;
  }
  try {
    const response = await fetch(`${API_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, username, password, isAdmin }),
    });
    const data = await response.json();
    if (!response.ok) {
      errorMessage.textContent = data.error || "Ошибка регистрации";
      return;
    }
    localStorage.setItem("token", data.token);
    localStorage.setItem("email", data.email);
    localStorage.setItem("username", data.username);
    localStorage.setItem("isAdmin", data.isAdmin);
    window.location.href = "chat.html";
  } catch (err) {
    errorMessage.textContent = "Ошибка соединения с сервером";
  }
});
