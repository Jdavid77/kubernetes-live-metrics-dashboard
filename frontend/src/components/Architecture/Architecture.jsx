import React from 'react';
import FlowDiagram from './FlowDiagram';
import config from '../../config/env';
import './Architecture.css';

// Topics come from the `architecture` section of the runtime config.yaml.
const Architecture = () => {
  const topics = config.ARCHITECTURE_TOPICS;

  if (!topics) {
    return <div className="text-center text-muted">Failed to load architecture</div>;
  }

  return (
    <div className="architecture">
      <div className="container">
        <nav className="architecture-toc">
          {topics.map((topic) => (
            <a key={topic.id} href={`#${topic.id}`} style={{ '--accent': topic.color }}>
              {topic.title}
            </a>
          ))}
        </nav>

        {topics.map((topic) => (
          <section
            key={topic.id}
            id={topic.id}
            className="topic-card"
            style={{ '--accent': topic.color }}
          >
            <h2 className="topic-title">{topic.title}</h2>
            <p className="topic-summary">{topic.summary}</p>

            {topic.lanes && <FlowDiagram lanes={topic.lanes} />}

            {topic.components && (
              <>
                <h3 className="topic-subtitle">Components</h3>
                <ul className="topic-components">
                  {topic.components.map((c) => (
                    <li key={c.name}>
                      <strong>{c.name}</strong>
                      <span>{c.role}</span>
                    </li>
                  ))}
                </ul>
              </>
            )}

            {topic.decisions && (
              <>
                <h3 className="topic-subtitle">Decisions</h3>
                <div className="topic-decisions">
                  {topic.decisions.map((d) => (
                    <div key={d.title} className="decision">
                      <h4>{d.title}</h4>
                      <p>{d.why}</p>
                    </div>
                  ))}
                </div>
              </>
            )}
          </section>
        ))}
      </div>
    </div>
  );
};

export default Architecture;
