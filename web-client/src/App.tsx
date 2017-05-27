import * as React from "react";
import Main from "./components/Main";
import Box from "./components/Box";
import DeviceList from "./components/DeviceList";
import Title from "./components/Title";
import Device from "./components/Device";
import InputField from "./components/form/InputField";
import Button from "./components/form/Button";
import Section from "./components/Section";

const App = () => (
  <div>
    <Main>
      <Box>
        <Title>Mielke/1.0</Title>
      </Box>
      <DeviceList>
        <Device online={true} alias="test" hostname="bla" />
        <Device online={false} alias="test2" hostname="blablablasdfasdlfjkhasdjfh" />
      </DeviceList>
      <Box>
        <Section flex={true}>
          <Button color="red">set invisible</Button>
          <InputField />
          <Button color="green">set visible</Button>
        </Section>
      </Box>
    </Main>
  </div>
);

export default App;
