import React from 'react';
import ReactMarkdown from 'react-markdown';

type Props = {
  markdown: string;
}

const Markdown: React.FC<Props> = ({markdown}) => {
  return (<ReactMarkdown>{markdown}</ReactMarkdown>);
}

export default Markdown;