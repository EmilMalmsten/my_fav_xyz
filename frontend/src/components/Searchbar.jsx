import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

function Searchbar() {
    const texts = [
        "Restaurants in London",
        "Cat cafes in Tokyo",
        "Nailbiting movies",
        "Travel themed books",
    ];
    const [index, setIndex] = useState(0);
    const [subIndex, setSubIndex] = useState(0);
    const [isErasing, setIsErasing] = useState(false);
    const [isInputFocused, setIsInputFocused] = useState(false);
    const [inputValue, setInputValue] = useState("");
    const navigate = useNavigate();

    useEffect(() => {
        if (!isInputFocused) {
            if (isErasing && subIndex > 0) {
                setTimeout(() => {
                    setSubIndex(subIndex - 1);
                }, 50);
            } else if (isErasing && subIndex === 0) {
                setIsErasing(false);
                if (index < texts.length - 1) {
                    setIndex(index + 1);
                } else {
                    setIndex(0);
                }
            } else if (subIndex < texts[index].length + 1) {
                setTimeout(() => {
                    setSubIndex(subIndex + 1);
                }, 150);
            } else if (subIndex === texts[index].length + 1) {
                setTimeout(() => {
                    setIsErasing(true);
                }, 1000);
            }
        }
    }, [subIndex, index, isErasing, isInputFocused, texts]);

    const handleInputFocus = () => {
        setIsInputFocused(true);
    };

    const handleInputBlur = () => {
        setIsInputFocused(false);
        setInputValue("");
    };

    const handleInputChange = (e) => {
        setInputValue(e.target.value);
    };

    const handleSearch = async (e) => {
        if (e.key === "Enter") {
            const limit = 20;
            const offset = 0;
            navigate(
                `/toplists/search?searchTerm=${inputValue}&limit=${limit}&offset=${offset}`
            );
        }
    };

    return (
        <div className="my-4 text-center">
            <h1 className="mb-0">Find a toplist for</h1>
            <div className="d-flex align-items-center justify-content-center">
                <div className="mx-1">🔍</div>
                {isInputFocused ? (
                    <input
                        className="dynamic-text-input"
                        value={inputValue}
                        placeholder="Enter search term"
                        onFocus={handleInputFocus}
                        onBlur={handleInputBlur}
                        onChange={handleInputChange}
                        onKeyDown={handleSearch}
                    />
                ) : (
                    <input
                        className="dynamic-text"
                        value={texts[index].substring(0, subIndex)}
                        readOnly
                        onFocus={handleInputFocus}
                    />
                )}
            </div>
        </div>
    );
}

export default Searchbar;
