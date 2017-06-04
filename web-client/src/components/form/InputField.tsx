import * as React from "react";
import styled from "styled-components";

export default styled.input`
  padding: 1rem 3rem;
  font-size: .8rem;
  border-radius: .25rem;
  outline: 0;
  border: 0;
  flex: 1;
  text-align: center;
  @media (max-width: 48rem) {
    flex: 0 0 auto;
  }
`;
