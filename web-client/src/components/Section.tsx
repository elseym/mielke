import * as React from "react";
import styled from "styled-components";

interface SectionProps {
  flex: boolean;
  children: JSX.Element[] | JSX.Element | string;
  className?: string;
}

const Section = ({children, className}: SectionProps) => (
  <div className={className}>{children}</div>
);

export default styled.div`
  border-radius: .25rem;
  outline: 0;
  border: 0;
  padding: .75rem;
  ${({flex}: SectionProps) => flex ? "display: flex;" : ""}
`;
