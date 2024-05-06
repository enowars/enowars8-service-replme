import { addLogin } from "./lib";

const form = document.querySelector("#new_container_form");

async function sendData() {
  // Associate the FormData object with the form element
  const formData = new FormData(form as HTMLFormElement);

  try {
    const response = await fetch("/api/term/private", {
      method: "POST",
      body: formData,
    });

    const payload = await response.json()

    if (response.ok) {
      const username = formData.get('username') as string;
      const password = formData.get('password') as string;
      addLogin(username, password)
      location.href = `/term/${username}/${payload.port}`;
    } else {
      document.getElementById("form_error_content")!.innerText = payload.error;
      document.getElementById("form_error_container")?.classList.remove("hidden");
      setTimeout(() => {
        document.getElementById("form_error_container")?.classList.add("hidden");
      }, 5000);
    }
  } catch (e) {
    console.error(e);
  }
}

form!.addEventListener("submit", (event) => {
  event.preventDefault();
  sendData();
});

