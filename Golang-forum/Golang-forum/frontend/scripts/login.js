const API_URL = "http://localhost:8080/api";
const emailInput = document.getElementById("email");
const passwordInput = document.getElementById("password");
const loginBtn = document.getElementById("login-btn");
const errorMessage = document.getElementById("error-message");

function validateEmail(email) {
  return /\S+@\S+\.\S+/.test(email);
}

loginBtn.addEventListener("click", async () => {
  const email = emailInput.value.trim();
  const password = passwordInput.value;
  if (!validateEmail(email) || !password) {
    errorMessage.textContent = "Заполните все поля корректно";
    return;
  }
  try {
    const response = await fetch(`${API_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    // Логируем статус ответа
    console.log("Response status:", response.status);

    const data = await response.json();

    // Логируем данные ответа
    console.log("Response data:", data);

    if (!response.ok) {
      errorMessage.textContent = data.error || "Ошибка входа";
      return;
    }
    localStorage.setItem("token", data.token);
    localStorage.setItem("email", data.email);
    localStorage.setItem("username", data.username);
    localStorage.setItem("isAdmin", data.isAdmin);
    window.location.href = "chat.html";
  } catch (err) {
    console.error("Error:", err);
    errorMessage.textContent = "Ошибка соединения с сервером";
  }
});
