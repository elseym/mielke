import * as React from "react";
import styled from "styled-components";

interface EllipsisProps {
  bold?: boolean;
  className?: string;
  children: JSX.Element | string;
}

const Ellipsis = ({bold, className, children}: EllipsisProps) => (
  <p className={className}>
    {children}
  </p>
);

export default styled(Ellipsis)`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
  ${({bold}: any) => bold ? "font-weight: 900;" : ""}
`;
