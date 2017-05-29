import * as React from "react";
import styled from "styled-components";
import Ellipsis from "./Ellipsis";

interface AvatarProps {
  url: string;
  className?: string;
  hostname: string;
  alias: string;
}

const Wrapper = styled.div`
  :after {
    content: "";
    display: table;
    clear: both;
  }
  border-radius: 2.1rem;
`;


const Name = styled.div`
  float: left;
  margin-left: .8rem;
  width: 65%;
`;

const Avatar = ({url, className, hostname, alias}: AvatarProps) => (
    <Wrapper>
      <div className={className} style={{float: "left"}} />
      <Name>
        <Ellipsis bold={true}>{alias}</Ellipsis>
        <Ellipsis>{hostname}</Ellipsis>
      </Name>
    </Wrapper>
);

export default styled(Avatar)`
  width: 4.2rem;
  height: 4.2rem;
  border-radius: 2.1rem;
  float: left;
  background-image: url(${ ({url}: AvatarProps) => url }?s=128&d=retro);
  background-size: cover;
`;
