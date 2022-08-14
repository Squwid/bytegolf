import React from 'react';
import Icon from './loading-icon.gif';

type Props = {
  style?: React.CSSProperties;
};

const LoadingIcon: React.FC<Props> = (props) => (<img style={props.style} src={Icon} alt="Loading Icon"/>);

export default LoadingIcon;