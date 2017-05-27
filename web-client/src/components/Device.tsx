import * as React from "react";
import styled, {StyledProps} from "styled-components";
import Ellipsis from "./Ellipsis";

interface DeviceProps {
  online: boolean;
  alias: string;
  hostname: string;
  className?: string;
}

const Device = ({online, alias, hostname, className}: DeviceProps) => (
  <div className={className}>
    <Ellipsis bold={online}>{alias}</Ellipsis>
    <Ellipsis bold={false}><small>{hostname}</small></Ellipsis>
  </div>
);

export default styled(Device)`
  width: 100%;
  @media (min-width: 48em) {
    width: 20%;
  }
  margin: 0.3rem;
  padding: .25rem .75rem;
  border-radius: .25rem;
  background: ${ ({online}: DeviceProps) => online ? "#e9ffd9" : "#ffecec" };
  color: ${ ({online}: DeviceProps) => online ? "#4f5459" : "#f5aca6" };
  ${({online}: DeviceProps) => online ? "" : "opacity: .4;"}
`;
