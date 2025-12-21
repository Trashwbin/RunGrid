import type {ChangeEvent, Ref} from 'react';
import {Icon} from '../ui/Icon';

type SearchBarProps = {
  value: string;
  onChange: (value: string) => void;
  inputRef?: Ref<HTMLInputElement>;
};

export function SearchBar({value, onChange, inputRef}: SearchBarProps) {
  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange(event.target.value);
  };

  return (
    <div className="search-bar">
      <input
        ref={inputRef}
        type="search"
        placeholder="快速搜索"
        value={value}
        onChange={handleChange}
      />
      <button type="button" className="favorite-button" aria-label="收藏">
        <Icon name="heart" />
      </button>
    </div>
  );
}
