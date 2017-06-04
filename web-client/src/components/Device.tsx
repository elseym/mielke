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
  lastSeen: string;
}

const Device = ({avatarURL, online, alias, hostname, className, lastSeen}: DeviceProps) => (
  <div className={className}>
    <Avatar url={avatarURL} hostname={hostname} alias={alias} />
    <Ellipsis>{lastSeen}</Ellipsis>
  </div>
);

export default styled(Device)`
  display: flex;
  flex-direction: column;
  margin: .5rem;
  width: 15rem;
  padding: .5rem;
  border-radius: .25rem;
  @media (max-width: 36rem) {
      width: 98%;
  }
  background: ${ ({online}: DeviceProps) => online ? "#e9ffd9" : "#ffecec" };
  color: ${ ({online}: DeviceProps) => online ? "#4f5459" : "#f5aca6" };
  ${({online}: DeviceProps) => online ? "" : "opacity: .4;"}
`;
