export const AboutPageContent = () => {
  return (
    <main className="mx-auto flex max-w-screen-xl flex-1 flex-col p-4">
      <h2 className="mb-4 text-2xl font-semibold">About</h2>

      <section className="flex flex-col gap-6">
        <div>
          The GitHub Trending page was not enough,
          <span className="font-semibold"> so I built my own.</span>
        </div>

        <div>
          <p>
            The data is fetched hourly from the official GitHub API by my custom
            built data loader written in Go. The database in use is TimescaleDB,
            I aggregate the data in dayly buckets and calculate the star
            difference for the time periods. In order to achieve less than 100ms
            query times, I store the end result in materialized views, which are
            updated hourly. The frontend is using Tanstack Start + Tailwindcss,
            very basic.
          </p>
        </div>

        <div>
          <p>
            I am using the GitHub graphql "Search" endpoint to fetch all the
            repositories ordered by stars descending. In order to overcome the
            maximum of 1000 search results limitation, I decrease my star range
            query every 1000 repos.
          </p>
          <p>
            Example: "stars:200..500000" --&gt; "stars:200..100000" --&gt;
            "stars:200..20000" --&gt; ...
          </p>
          <p>
            Around 240k repositories are fetched in 21 minutes. I found some
            little tricks to speed up the fetching like for example, fetching
            all 10 pages (100 pageSize * 10 == 1000 results) at once,
            prefetching the next star ranges, or making 200 requests at once and
            bypassing the secondary limit somehow.
          </p>
        </div>

        <div>
          Thanks for reading this and if you have any feedback, bug reports or
          feature requests, please let me know in my{' '}
          <a
            href="https://github.com/glup3/trendingrepos"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-muted-foreground underline"
          >
            GitHub issues
          </a>
          .
        </div>
      </section>
    </main>
  )
}
