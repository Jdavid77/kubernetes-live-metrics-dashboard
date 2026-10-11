import React from 'react';
import FlowDiagram from './FlowDiagram';
import topics from './topics';
import './Architecture.css';

const Architecture = () => (
  <div className="architecture">
    <div className="container">
      <nav className="architecture-toc">
        {topics.map((topic) => (
          <a key={topic.id} href={`#${topic.id}`} style={{ '--accent': topic.color }}>{topic.title}</a>
        ))}
      </nav>

      {topics.map((topic) => (
        <section key={topic.id} id={topic.id} className="topic-card" style={{ '--accent': topic.color }}>
          <h2 className="topic-title">{topic.title}</h2>
          <p className="topic-summary">{topic.summary}</p>

          {topic.lanes && <FlowDiagram lanes={topic.lanes} />}

          <h3 className="topic-subtitle">Components</h3>
          <ul className="topic-components">
            {topic.components.map((c) => (
              <li key={c.name}>
                <strong>{c.name}</strong>
                <span>{c.role}</span>
              </li>
            ))}
          </ul>

          <h3 className="topic-subtitle">Decisions</h3>
          <div className="topic-decisions">
            {topic.decisions.map((d) => (
              <div key={d.title} className="decision">
                <h4>{d.title}</h4>
                <p>{d.why}</p>
              </div>
            ))}
          </div>
        </section>
      ))}
    </div>
  </div>
);

export default Architecture;
