import React from 'react';
import TOCOriginal from '@theme-original/TOC';

export default function TOCWrapper(props: any) {
  return (
    <div className="custom-toc-container">
      <TOCOriginal {...props} />
      <div className="toc-review-section">
        <hr className="toc-review-divider" />
        <h4 className="toc-review-title">Was this helpful?</h4>
        <div className="toc-review-stars">
          {[1, 2, 3, 4, 5].map((star) => (
            <svg 
              key={star} 
              className="star-icon" 
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              strokeWidth="1.5" 
              strokeLinecap="round" 
              strokeLinejoin="round"
            >
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
            </svg>
          ))}
        </div>
      </div>
    </div>
  );
}
