import * as React from "react";
import styled from "styled-components";
import {defaultFont} from "../styles/mixins";

interface ButtonProps {
  color: "green" | "red";
  children: string;
  className?: string;
  onClick?: (e: React.MouseEvent<HTMLButtonElement>) => void;
}

const Button = ({color, children, className, onClick}: ButtonProps) => (
  <button className={className} onClick={onClick}>
    {children}
  </button>
);

export default styled(Button)`
  padding: 1rem 3rem;
  font-size: .8rem;
  border-radius: .25rem;
  outline: 0;
  border: 0;
  background-color: ${({color}: ButtonProps) => color === "green" ? "#67ba72;" : "#ff7772"};
  text-transform: uppercase;
  text-align: center;
  flex: 1;
  ${defaultFont}
  @media (max-width: 48rem) {
    flex: 0 0 auto;
  }
`;
