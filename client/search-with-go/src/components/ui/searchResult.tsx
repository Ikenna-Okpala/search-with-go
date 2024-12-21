import { SearchResultType } from "@/Search";

const getBaseUrl = (url: string) => {
  if (url.length === 0) {
    return "";
  }
  const paths = url.split("/");

  return paths[0] + "//" + paths[2];
};

function SearchResult({
  icon,
  name,
  url,
  title,
  description,
}: SearchResultType) {
  return (
    <div className="flex flex-col w-3/5 gap-y-2">
      <div className="flex flex-row items-center gap-x-3">
        <img
          src={icon}
          width={35}
          height={35}
          className="rounded-full object-cover"
        />

        <div className="flex flex-col gap-y-0.5">
          <span className="text-white text-xl">{name}</span>
          <span className="text-white">{getBaseUrl(url)}</span>
        </div>
      </div>
      <a
        href={
          url
            ? url
            : "https://storage.googleapis.com/search-with-go/world-wide-web.png"
        }
        className="text-blue-400 text-2xl hover:underline w-fit"
      >
        {title ? title : "No title provided."}
      </a>
      <p className="text-slate-200 line-clamp-3 text-lg">
        {description ? description : "No description provided."}
      </p>
    </div>
  );
}

export default SearchResult;
