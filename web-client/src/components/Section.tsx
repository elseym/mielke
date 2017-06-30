import * as React from "react";
import styled from "styled-components";

interface SectionProps {
  flex: boolean;
  children: JSX.Element[] | JSX.Element | string;
  className?: string;
}

const Wrapper = styled.div`
  display: flex;
  padding: .5rem 1rem;
  @media (max-width: 64rem) {
    flex-direction: column;
  }
`;

const Section = ({children, className}: SectionProps) => (
  <section className={className}>
    {React.Children.map(children, (child, i) => {
      return <Wrapper>{child}</Wrapper>;
    })}
  </section>
);

export default styled(Section)`
  border-radius: .25rem;
  background-color: #e2e5e9;
  outline: 0;
  border: 0;
  padding: .75rem;
  ${({flex}: SectionProps) => {
    if (flex) {
      return `
        display: flex;
        justify-content: space-between;
        @media (max-width: 64rem) {
          flex-direction: column;
        }
      `;
    }

    return "";
  }}
`;
