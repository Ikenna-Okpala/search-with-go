import { IoSearch } from "react-icons/io5";
import Go from "./assets/go-mascot.svg";
import { Input } from "./components/ui/input";
import { Separator } from "./components/ui/separator";
import SearchResult from "./components/ui/searchResult";
import { useLocation } from "react-router-dom";
import { useEffect, useState } from "react";

export type SearchResultType = {
  url: string;
  title: string;
  description: string;
  name: string;
  icon: string;
};

function Search() {
  const { state } = useLocation();

  const search = state?.search;

  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<SearchResultType[]>([]);

  const fetchSearchResults = async (search: string) => {
    try {
      const url = new URL("http://localhost:8080/");

      const params = { query: search };

      url.search = new URLSearchParams(params).toString();

      const resp = await fetch(url);

      const searchResults: SearchResultType[] = await resp.json();

      console.log("results: ", searchResults);

      setSearchResults(searchResults);
    } catch (err) {
      console.log(err);
    }
  };

  useEffect(() => {
    if (search) {
      setSearchQuery(search);
      fetchSearchResults(search);
    }
  }, []);

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
  };

  const onEnterPressed = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      fetchSearchResults(searchQuery);
    }
  };

  return (
    <div className="w-screen h-screen flex-col overflow-y-auto">
      <div className="flex flex-row gap-x-4 px-20 pt-10 pb-8">
        <div className="flex flex-row gap-x-3 items-center">
          <span className="text-4xl font-medium text-white">Search With</span>

          <img src={Go} alt="Go mascot" width={50} height={50} />
        </div>

        <div className="flex flex-row gap-x-3 w-2/4 items-center bg-gray200 rounded-full px-5">
          <IoSearch className="text-4xl text-gray100" />
          <Input
            type="text"
            className="h-[60px] bg-transparent border-none text-slate-50 text-xl"
            value={searchQuery}
            onChange={onChange}
            onKeyDown={onEnterPressed}
          />
        </div>
      </div>

      <Separator orientation="horizontal" className="bg-gray100" />
      <div className="flex flex-col gap-y-6 px-20 py-10">
        {searchResults.length > 0 ? (
          searchResults.map((searchResult, index) => (
            <SearchResult
              name={searchResult.name}
              description={searchResult.description}
              icon={searchResult.icon}
              title={searchResult.title}
              url={searchResult.url}
              key={index}
            />
          ))
        ) : (
          <div>
            <ul className="marker:text-white list-disc space-y-3">
              <li className="text-slate-50 text-lg">
                No result matched your search
              </li>
              <li className="text-white text-lg">
                Ensure your words are spelled correctly
              </li>
              <li className="text-white text-lg">Try more generic keywords</li>
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}

export default Search;
