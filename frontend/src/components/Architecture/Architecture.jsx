import React, { useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import FlowDiagram from './FlowDiagram';
import { topicIcon } from './icons';
import config from '../../config/env';
import './Architecture.css';

const DEFAULT_ACCENT = 'var(--color-primary)';

const plural = (n, word) => `${n} ${word}${n === 1 ? '' : 's'}`;

// Topics come from the `architecture` section of the runtime config.yaml.
// The active topic lives in the URL hash, e.g. /architecture#storage.
const Architecture = () => {
  const topics = config.ARCHITECTURE_TOPICS;
  const { hash } = useLocation();

  const active = topics && (topics.find((t) => `#${t.id}` === hash) || topics[0]);

  useEffect(() => {
    window.scrollTo({ top: 0 });
  }, [hash]);

  if (!topics) {
    return <div className="text-center text-muted">Failed to load architecture</div>;
  }

  const ActiveIcon = topicIcon(active.icon);
  const meta = [
    active.components?.length && plural(active.components.length, 'component'),
    active.decisions?.length && plural(active.decisions.length, 'decision'),
    active.lanes?.length && plural(active.lanes.length, 'flow'),
  ].filter(Boolean);

  return (
    <div className="architecture" style={{ '--accent': active.color || DEFAULT_ACCENT }}>
      <div className="container architecture-layout">
        <nav className="architecture-nav" aria-label="Architecture topics">
          {topics.map((topic) => {
            const Icon = topicIcon(topic.icon);
            const isActive = topic.id === active.id;
            return (
              <Link
                key={topic.id}
                to={`#${topic.id}`}
                className={`architecture-nav-item${isActive ? ' is-active' : ''}`}
                style={{ '--item-accent': topic.color || DEFAULT_ACCENT }}
                aria-current={isActive ? 'page' : undefined}
              >
                <Icon aria-hidden="true" />
                <span>{topic.title}</span>
              </Link>
            );
          })}
        </nav>

        <article key={active.id} className="topic">
          <header className="topic-hero">
            <div className="topic-icon">
              <ActiveIcon aria-hidden="true" />
            </div>
            <div>
              <h2 className="topic-title">{active.title}</h2>
              <p className="topic-summary">{active.summary}</p>
              {meta.length > 0 && <p className="topic-meta">{meta.join(' · ')}</p>}
            </div>
          </header>

          {active.lanes && <FlowDiagram lanes={active.lanes} />}

          {active.components && (
            <section>
              <h3 className="topic-subtitle">Components</h3>
              <ul className="topic-components">
                {active.components.map((c) => (
                  <li key={c.name}>
                    <strong>{c.name}</strong>
                    <span>{c.role}</span>
                  </li>
                ))}
              </ul>
            </section>
          )}

          {active.decisions && (
            <section>
              <h3 className="topic-subtitle">Decisions</h3>
              <ol className="topic-decisions">
                {active.decisions.map((d, i) => (
                  <li key={d.title}>
                    <span className="decision-number" aria-hidden="true">
                      {String(i + 1).padStart(2, '0')}
                    </span>
                    <div>
                      <h4>{d.title}</h4>
                      <p>{d.why}</p>
                    </div>
                  </li>
                ))}
              </ol>
            </section>
          )}
        </article>
      </div>
    </div>
  );
};

export default Architecture;
