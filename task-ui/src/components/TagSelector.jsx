import { useState, useRef, useEffect } from 'react';
import './TagSelector.css';

// 23 predefined tags
const PROFESSION_TAGS = {
  // 技术类
  'frontend-developer': '前端开发',
  'backend-developer': '后端开发',
  'fullstack-developer': '全栈开发',
  'mobile-developer': '移动开发',
  'devops-engineer': 'DevOps 工程师',
  'data-engineer': '数据工程师',
  'ml-engineer': '机器学习工程师',
  'qa-engineer': '测试工程师',
  
  // 设计类
  'ui-designer': 'UI 设计师',
  'ux-designer': 'UX 设计师',
  'product-designer': '产品设计师',
  'graphic-designer': '平面设计师',
  
  // 产品/管理类
  'product-manager': '产品经理',
  'project-manager': '项目经理',
  'scrum-master': '敏捷教练',
  
  // 商业类
  'business-analyst': '商业分析师',
  'entrepreneur': '创业者',
  'consultant': '咨询顾问',
  
  // 其他
  'researcher': '研究员',
  'writer': '写作者',
  'marketer': '市场营销',
};

function TagSelector({ selectedTags, onChange, maxTags = 5, disabled = false }) {
  const [inputValue, setInputValue] = useState('');
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [filteredTags, setFilteredTags] = useState([]);
  const inputRef = useRef(null);
  const suggestionsRef = useRef(null);

  // Filter predefined tags based on input
  useEffect(() => {
    if (inputValue.length > 0) {
      const filtered = Object.entries(PROFESSION_TAGS).filter(([key, label]) => {
        const searchTerm = inputValue.toLowerCase();
        return (
          key.toLowerCase().includes(searchTerm) ||
          label.toLowerCase().includes(searchTerm)
        );
      });
      setFilteredTags(filtered);
      setShowSuggestions(filtered.length > 0);
    } else {
      setFilteredTags([]);
      setShowSuggestions(false);
    }
  }, [inputValue]);

  // Close suggestions when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        suggestionsRef.current &&
        !suggestionsRef.current.contains(event.target) &&
        !inputRef.current.contains(event.target)
      ) {
        setShowSuggestions(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleAddTag = (tag) => {
    if (selectedTags.length >= maxTags) {
      return;
    }
    if (!selectedTags.includes(tag)) {
      onChange([...selectedTags, tag]);
    }
    setInputValue('');
    setShowSuggestions(false);
  };

  const handleRemoveTag = (tag) => {
    onChange(selectedTags.filter((t) => t !== tag));
  };

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleInputKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      const trimmedValue = inputValue.trim();
      
      // If there's a matching predefined tag, add it
      if (filteredTags.length > 0) {
        handleAddTag(filteredTags[0][0]);
      } else if (trimmedValue) {
        // Add as custom tag
        if (trimmedValue.length <= 50) {
          handleAddTag(trimmedValue);
        }
      }
    }
  };

  const handleSelectPredefined = (key) => {
    handleAddTag(key);
  };

  const handleAddCustom = () => {
    const trimmedValue = inputValue.trim();
    if (trimmedValue && trimmedValue.length <= 50) {
      handleAddTag(trimmedValue);
    }
  };

  const isPreDefinedTag = (tag) => {
    return Object.keys(PROFESSION_TAGS).includes(tag);
  };

  return (
    <div className="tag-selector">
      {/* Selected tags display */}
      {selectedTags.length > 0 && (
        <div className="selected-tags">
          {selectedTags.map((tag) => (
            <div
              key={tag}
              className={`selected-tag ${isPreDefinedTag(tag) ? 'predefined' : 'custom'}`}
            >
              <span className="tag-text">
                {isPreDefinedTag(tag) ? PROFESSION_TAGS[tag] : tag}
              </span>
              <button
                type="button"
                className="remove-tag-btn"
                onClick={() => handleRemoveTag(tag)}
                disabled={disabled}
                title="移除标签"
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Tag input */}
      <div className="tag-input-container">
        <input
          ref={inputRef}
          type="text"
          className="tag-input"
          value={inputValue}
          onChange={handleInputChange}
          onKeyDown={handleInputKeyDown}
          placeholder={
            selectedTags.length >= maxTags
              ? `最多 ${maxTags} 个标签`
              : '自定义职业标签...'
          }
          disabled={disabled || selectedTags.length >= maxTags}
        />
        {inputValue.trim() && selectedTags.length < maxTags && (
          <button
            type="button"
            className="add-custom-btn"
            onClick={handleAddCustom}
            disabled={disabled}
          >
            添加
          </button>
        )}
      </div>

      {/* Suggestions dropdown */}
      {showSuggestions && (
        <div ref={suggestionsRef} className="tag-suggestions">
          {filteredTags.map(([key, label]) => (
            <div
              key={key}
              className="tag-suggestion-item"
              onClick={() => handleSelectPredefined(key)}
            >
              <span className="suggestion-label">{label}</span>
              <span className="suggestion-key">{key}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default TagSelector;
