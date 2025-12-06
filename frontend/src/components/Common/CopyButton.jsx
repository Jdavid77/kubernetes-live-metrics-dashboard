import React, { useState } from 'react';
import { FiCopy, FiCheck } from 'react-icons/fi';
import './CopyButton.css';

const CopyButton = ({ text }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  return (
    <button
      className={`copy-button ${copied ? 'copied' : ''}`}
      onClick={handleCopy}
      title="Copy to clipboard"
    >
      {copied ? <FiCheck /> : <FiCopy />}
    </button>
  );
};

export default CopyButton;
