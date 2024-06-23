import { randomString, UserSessionResponseSchema } from "./lib"

customElements.define(
  'session-elem',
  class extends HTMLElement {
    constructor() {
      super();
      // @ts-ignore
      const template = document.getElementById("session-template").content;
      const shadowRoot = this.attachShadow({ mode: "open" });
      [
        '/static/css/style.css',
        '/static/fontawesome/css/fontawesome.css',
        '/static/fontawesome/css/brands.css',
        '/static/fontawesome/css/solid.css'
      ].forEach((href) => {
        const linkElem = document.createElement('link');
        linkElem.rel = 'stylesheet';
        linkElem.href = href;
        shadowRoot.appendChild(linkElem);
      })
      shadowRoot.appendChild(template.cloneNode(true));
    }
    connectedCallback() {
      const hrefSlot = this.shadowRoot.querySelector('slot[name="href"]');
      const dynamicLink = this.shadowRoot.getElementById('dynamic-link');

      const updateHref = () => {
        // @ts-ignore
        const nodes = hrefSlot.assignedNodes();
        if (nodes.length > 0) {
          // @ts-ignore
          dynamicLink.href = nodes[0].textContent;
        }
      };

      hrefSlot.addEventListener('slotchange', updateHref);
      updateHref();
    }
  },
)

const button = document.getElementById("repl-button")

button.onclick = async () => {
  const username = randomString(60);
  const password = randomString(60);
  const res = await fetch(
    '/api/repl',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username,
        password
      })
    }
  )

  if (res.ok) {
    const id = (await res.json()).id
    location.href = "/term/" + id
  }

}

fetch("/api/user/sessions")
  .then((response) => response.json())
  .then((data) => UserSessionResponseSchema.parse(data))
  .then((names) => {
    const container = document.getElementById("session-container")
    names.forEach((name) => {

      const session = document.createElement('session-elem');
      const labelSlot = document.createElement('slot');
      labelSlot.slot = "label";
      labelSlot.innerText = name.substring(0, 10) + "...";
      const hrefSlot = document.createElement('slot');
      hrefSlot.slot = "href";
      hrefSlot.innerText = "/term/" + name;
      session.appendChild(labelSlot)
      session.appendChild(hrefSlot)

      container.appendChild(session);
    })
  })

