import { UserSessionResponseSchema } from "./lib"

fetch("/api/user/sessions")
  .then((response) => response.json())
  .then((data) => UserSessionResponseSchema.parse(data))
  .then((names) => {
    const container = document.getElementById("session_container")
    names.forEach((name) => {

      const row = document.createElement('div')
      row.classList.add(
        'flex',
        'flex-row',
        'w-full',
        'justify-between'
      )

      const label = document.createElement('div')
      label.innerText = name.substring(0, 10);

      const link = document.createElement('a')
      link.href = "/term/" + name
      link.classList.add(
        "rounded",
        "bg-indigo-600",
        "px-2",
        "py-1",
        "font-semibold",
        "text-white",
        "shadow-sm",
        "hover:bg-indigo-500",
        "focus-visible:outline",
        "focus-visible:outline-2",
        "focus-visible:outline-offset-2",
        "focus-visible:outline-indigo-600"
      )

      const icon = document.createElement('i')
      icon.classList.add(
        'fa-solid',
        'fa-arrow-right'
      )

      link.appendChild(icon);

      row.appendChild(label);
      row.appendChild(link);

      container.appendChild(row)
    })
  })


fetch("/api/user/sessions/debug")
  .then((response) => response.json())
  .then((data) => {
    const container = document.getElementById("sessions_debug")
    container.innerText = JSON.stringify(data, null, 2)
  })

