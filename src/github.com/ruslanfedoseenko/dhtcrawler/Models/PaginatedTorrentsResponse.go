package Models

type PaginatedTorrentsResponse struct {
	Page         int
	PageCount    int
	ItemsCount   int
	ItemsPerPage int
	Torrents     []Torrent
}
