import * as React from "react";
import styled from "styled-components";
import {defaultFont} from "./styles/mixins";

export default styled.main`
  width: 48rem;
  margin: 0 auto;
  ${defaultFont}
  @media (max-width: 48rem) {
    width: 98%;
  }
`;
