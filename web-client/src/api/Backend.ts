export interface BackendData {
  self: BackendClient;
  list: {
    [key: string]: BackendClient;
  }
}

export interface BackendClient {
  alias: string;
  hostname: string;
  online: boolean;
  associated: string;
  lastSeen: string;
  ap: string;
}

export default class Backend {
  getData() {
    const headers = new Headers();
    headers.set("Accept", "application/json");

    return fetch("?", {
      headers,
    });
  }

  setInvisible() {
    return fetch("?", {
      method: "DELETE",
    });
  }

  setVisible(alias: string) {
    const headers = new Headers();
    headers.set("Content-Type", "application/json");

    return fetch("?", {
      method: "PUT",
      headers,
      body: JSON.stringify({
        alias,
      }),
    });
  }
}
