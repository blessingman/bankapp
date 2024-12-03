const API_URL = "http://localhost:8080";

// Sections
const sections = {
  login: document.getElementById("login-section"),
  register: document.getElementById("register-section"),
  accounts: document.getElementById("accounts-section"),
};

// Buttons
const buttons = {
  login: document.getElementById("login-button"),
  register: document.getElementById("register-button"),
  addAccount: document.getElementById("add-account-button"),
  transfer: document.getElementById("transfer-button"),
  logout: document.getElementById("logout-button"),
};

// Error messages
const errors = {
  login: document.getElementById("login-error"),
  register: document.getElementById("register-error"),
};

// Inputs
const inputs = {
  username: document.getElementById("login-username"),
  password: document.getElementById("login-password"),
  newAccountBalance: document.getElementById("new-account-balance"),
  fromAccountId: document.getElementById("from-account-id"),
  toAccountId: document.getElementById("to-account-id"),
  transferAmount: document.getElementById("transfer-amount"),
};

// Message Container
const messageContainer = document.createElement("div");
messageContainer.id = "message-container";
document.body.appendChild(messageContainer);

messageContainer.style.cssText = `
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  background-color: #f0ad4e;
  color: white;
  border-radius: 5px;
  z-index: 1000;
  display: none;
`;

function showMessage(text, duration = 3000) {
  messageContainer.textContent = text;
  messageContainer.style.display = "block";

  setTimeout(() => {
    messageContainer.style.display = "none";
  }, duration);
}

// General Initialization
document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".container").forEach(container => {
    container.style.opacity = "0";
    container.style.transform = "translateY(20px)";
    container.style.transition = "opacity 0.6s ease, transform 0.6s ease";
    setTimeout(() => {
      container.style.opacity = "1";
      container.style.transform = "translateY(0)";
    }, 100);
  });
});

// Utility: Reset Forms and Focus
function resetFormsAndFocus(targetInput = null) {
  document.querySelectorAll("input").forEach(input => {
    if (input !== targetInput) {
      input.value = "";
    }
  });

  setTimeout(() => {
    if (targetInput) {
      targetInput.focus();
    }
  }, 50);
}

// Utility: Switch Sections
function switchSection(from, to, targetInput = null) {
  from.classList.add("hidden");
  to.classList.remove("hidden");

  if (targetInput) {
    setTimeout(() => targetInput.focus(), 100);
  }
}

// Login
buttons.login.addEventListener("click", async () => {
  const username = inputs.username.value;
  const password = inputs.password.value;

  try {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();
    if (!data.token) throw new Error("Ошибка авторизации");

    localStorage.setItem("token", data.token);
    switchSection(sections.login, sections.accounts, inputs.newAccountBalance);
    await fetchAndDisplayAccounts();
  } catch (error) {
    errors.login.textContent = error.message;
  } finally {
    resetFormsAndFocus(inputs.username);
  }
});

// Register
buttons.register.addEventListener("click", async () => {
  const username = document.getElementById("register-username").value;
  const password = document.getElementById("register-password").value;

  try {
    const response = await fetch(`${API_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();
    if (!data.success) throw new Error("Ошибка регистрации");

    showMessage("Регистрация успешна! Войдите в систему.", 3000);
    switchSection(sections.register, sections.login, inputs.username);
  } catch (error) {
    errors.register.textContent = error.message;
  } finally {
    resetFormsAndFocus(inputs.username);
  }
});

// Add Account
buttons.addAccount.addEventListener("click", async () => {
  const balance = parseFloat(inputs.newAccountBalance.value);

  if (isNaN(balance) || balance <= 0) {
    showMessage("Введите корректное числовое значение для баланса!", 3000);
    return;
  }

  try {
    await fetchWithAuth("/accounts", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ initial_balance: balance }),
    });
    showMessage("Счёт успешно добавлен!", 3000);
    await fetchAndDisplayAccounts();
  } catch (error) {
    showMessage(`Ошибка: ${error.message}`, 3000);
  } finally {
    resetFormsAndFocus(inputs.newAccountBalance);
  }
});

// Transfer Funds
buttons.transfer.addEventListener("click", async () => {
  const fromAccount = parseInt(inputs.fromAccountId.value, 10);
  const toAccount = parseInt(inputs.toAccountId.value, 10);
  const amount = parseFloat(inputs.transferAmount.value);

  if (isNaN(fromAccount) || isNaN(toAccount) || isNaN(amount) || amount <= 0) {
    showMessage("Введите корректные значения для перевода.", 3000);
    return;
  }

  try {
    await fetchWithAuth("/transfer", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        from_account_id: fromAccount,
        to_account_id: toAccount,
        amount: amount,
      }),
    });
    showMessage("Перевод выполнен успешно!", 3000);
    await fetchAndDisplayAccounts();
  } catch (error) {
    showMessage("Ошибка при выполнении перевода.", 3000);
  } finally {
    resetFormsAndFocus(inputs.transferAmount);
  }
});

// Logout
buttons.logout.addEventListener("click", () => {
  localStorage.removeItem("token");
  switchSection(sections.accounts, sections.login, inputs.username);
});
