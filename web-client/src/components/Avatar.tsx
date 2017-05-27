import * as React from "react";
import styled from "styled-components";

interface AvatarProps {
  url: string;
  className?: string;
}

const Avatar = ({url, className}: AvatarProps) => (
    <img className={className} />
);

export default styled(Avatar)`
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: url(${ ({url}: AvatarProps) => url }?s=128&d=retro) no-repeat 0 0;
`;
