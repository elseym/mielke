import * as React from "react";
import Main from "./components/Main";
import Box from "./components/Box";
import DeviceList from "./components/DeviceList";
import Title from "./components/Title";
import Device from "./components/Device";
import InputField from "./components/form/InputField";
import Button from "./components/form/Button";
import Section from "./components/Section";
import Backend from "./api/Backend";

const backend = new Backend();

interface AppState {
  data?: any;
  alias?: any
}

class App extends React.Component<any, AppState> {
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
    backend.getData().then((response: any) => {
      return response.json();
    }).then((data: any) => {
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
      return list.map((client: any, index: any) =>
        <Device online={client.online} alias={client.alias} hostname={client.hostname} />
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

  handleInputChange(e: any) {
    const alias = e.target.value;

    this.setState({
      alias,
    });
  }

  handleInvisible(e: any) {
    e.preventDefault();
    backend.setInvisible().then(() => {
      this.refresh();
    });
  }

  handleVisible(e: any) {
    e.preventDefault();
    backend.setVisible(this.state.alias).then(() => {
      this.refresh();
    });
  }
}

export default App;
