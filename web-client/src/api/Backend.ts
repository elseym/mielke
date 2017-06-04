export interface BackendClient {
  alias: string;
  avatarURL: string;
  hostname: string;
  online: boolean;
  associated: string;
  lastSeen: string;
  ap: string;
}

export interface BackendData {
  self: BackendClient;
  list: BackendClient[],
}

export default class Backend {
  getData(): Promise<BackendData> | undefined {
    const headers = new Headers();
    headers.set("Accept", "application/json");

    return fetch("?", {
      headers,
    })
      .then((resp: Response) => resp.json() as Promise<BackendData>)
      .then((data: BackendData) => data); // TODO: Verify unsafeData value here
  }

  static setInvisible() {
    return fetch("?", {
      method: "DELETE",
    });
  }

  static setVisible(alias: string) {
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
