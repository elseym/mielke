import * as React from "react";
import Main from "./components/Main";
import Box from "./components/Box";
import DeviceList from "./components/DeviceList";
import Title from "./components/Title";
import Device from "./components/Device";
import InputField from "./components/form/InputField";
import Button from "./components/form/Button";
import Section from "./components/Section";
import Backend, {BackendData, BackendClient} from "./api/Backend";

const backend = new Backend();

interface AppState {
  data?: BackendData;
  alias?: string;
}

class App extends React.Component<{}, AppState> {
  state: AppState = {};

  constructor() {
    super();

    this.state = {};

    this.renderList = this.renderList.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
    this.handleInvisible = this.handleInvisible.bind(this);
    this.handleVisible = this.handleVisible.bind(this);
    this.refresh = this.refresh.bind(this);
  }

  componentDidMount() {
    this.refresh();
  }

  refresh() {
    backend.getData().then((data: BackendData) => {
      this.setState({
        data,
        alias: data.self.alias,
      });
    });
  }

  renderList(): JSX.Element[] | null {
    const data = this.state.data;

    if (data !== undefined) {
      const list = data.list;
      return list.map((client: BackendClient, index: number) =>
        <Device
          avatarURL={client.avatarURL}
          online={client.online}
          alias={client.alias}
          hostname={client.hostname}
          lastSeen={client.lastSeen}
        />
      );
    }

    return null;
  }

  render() {
    const alias = this.state.data && this.state.data.self.alias;

    return (
      <div>
        <Main>
          <Box>
            <Title>Mielke/1.0</Title>
          </Box>
          <DeviceList>
            {this.renderList()}
          </DeviceList>
          <Section flex={true}>
            <Button color="red" onClick={this.handleInvisible}>set invisible</Button>
            <InputField onChange={this.handleInputChange} placeholder={alias} />
            <Button color="green" onClick={this.handleVisible}>set visible</Button>
          </Section>
        </Main>
      </div>
    )
  }

  handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    const alias = e.target.value;

    this.setState({
      alias,
    });
  }

  handleInvisible(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    Backend.setInvisible().then(() => {
      this.refresh();
    });
  }

  handleVisible(e: React.MouseEvent<HTMLButtonElement>) {
    e.preventDefault();
    Backend.setVisible(this.state.alias).then(() => {
      this.refresh();
    });
  }
}

export default App;
