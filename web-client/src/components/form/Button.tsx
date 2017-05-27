import * as React from "react";
import styled from "styled-components";

interface ButtonProps {
  color: "green" | "red";
  children: string;
  className?: string;
}

const Button = ({color, children, className}: ButtonProps) => (
  <button className={className}>
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
`;
