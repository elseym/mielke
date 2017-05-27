import * as React from "react";
import styled, {StyledProps} from "styled-components";
import Avatar from "./Avatar";
import Ellipsis from "./Ellipsis";

interface DeviceProps {
  avatarURL: string;
  online: boolean;
  alias: string;
  hostname: string;
  className?: string;
}

const Device = ({avatarURL, online, alias, hostname, className}: DeviceProps) => (
  <div className={className}>
    <Avatar url={avatarURL} />
    <Ellipsis bold={online}>{alias}</Ellipsis><br />
    <Ellipsis bold={false}><small>{hostname}</small></Ellipsis>
  </div>
);

export default styled(Device)`
  display: flex;
  flex-direction: row;
  margin: 0.3rem;
  padding: .25rem .75rem;
  border-radius: .25rem;
  background: ${ ({online}: DeviceProps) => online ? "#e9ffd9" : "#ffecec" };
  color: ${ ({online}: DeviceProps) => online ? "#4f5459" : "#f5aca6" };
  ${({online}: DeviceProps) => online ? "" : "opacity: .4;"}
`;
