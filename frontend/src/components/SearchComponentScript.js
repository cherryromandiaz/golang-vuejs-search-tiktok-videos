export default {
  props: {
    maxResults: { type: Number, default: 12 },
    buttonColor: { type: String, default: '#1da1f2' },
    placeholderText: { type: String, default: 'Search for TikTok videos' },
  },
  data() {
    return {
      query: '',
      videos: [],
      embedHtml: '',
      selectedVideoIndex: null,
      sortOption: 'relevance',
      fetchError: false,
    };
  },
  computed: {
    sortedVideos() {
      return this.sortOption === 'date'
        ? [...this.limitedVideos].sort((a, b) => this.parseTikTokDate(b.age) - this.parseTikTokDate(a.age))
        : this.limitedVideos;
    },
    limitedVideos() {
      return this.videos.slice(0, this.maxResults); // Limit by maxResults prop
    },
  },
  methods: {
    async fetchVideos() {
      if (this.query.trim() === '') {
        this.videos = [];
        this.fetchError = false;
        return;
      }

      // Check cache before making a request
      const cachedData = JSON.parse(localStorage.getItem(`cache_${this.query}`));
      const cacheExpiry = 10 * 60 * 1000; // Set cache expiry to 10 minutes
      const currentTime = new Date().getTime();

      if (cachedData && currentTime - cachedData.timestamp < cacheExpiry) {
        this.videos = cachedData.videos;
        this.fetchError = false;
        this.$emit('videosFetched', this.videos); // Emit event
        return;
      }

      try {
        const response = await fetch(`http://localhost:8080/search?query=${this.query}`);
        if (!response.ok) throw new Error('Failed to fetch videos');

        const data = await response.json();
        this.fetchError = false;

        // Map video data and add thumbnails
        this.videos = await Promise.all(
          data.map(async (video) => {
            try {
              const embedResponse = await fetch(`https://www.tiktok.com/oembed?url=${encodeURIComponent(video.url)}`);
              const embedData = await embedResponse.json();
              return {
                url: video.url,
                title: embedData.title ? embedData.title.slice(0, 100) + (embedData.title.length > 100 ? '...' : '') : "No Title",
                imageUrl: embedData.thumbnail_url,
                age: video.age,
              };
            } catch (error) {
              console.error('Error fetching thumbnail:', error);
              return { url: video.url, title: 'No Title', imageUrl: '', age: '' };
            }
          })
        );

        // Cache the new data with a timestamp
        localStorage.setItem(
          `cache_${this.query}`,
          JSON.stringify({ videos: this.videos, timestamp: currentTime })
        );
        this.$emit('videosFetched', this.videos); // Emit event
      } catch (error) {
        console.error('Error fetching videos:', error);
        this.fetchError = true;
      }
    },
    async fetchEmbed(videoUrl, index) {
      try {
        const response = await fetch(`https://www.tiktok.com/oembed?url=${encodeURIComponent(videoUrl)}`);
        if (!response.ok) throw new Error('Failed to fetch embed');

        const data = await response.json();
        this.embedHtml = data.html;
        this.selectedVideoIndex = index;
        this.$nextTick(() => {
          const script = document.createElement('script');
          script.src = 'https://www.tiktok.com/embed.js';
          script.async = true;
          document.body.appendChild(script);
          this.$emit('embedLoaded', index); // Emit event
        });
      } catch (error) {
        console.error('Error fetching embed:', error);
      }
    },
    parseTikTokDate(dateStr) {
      const currentYear = new Date().getFullYear();
      if (dateStr.match(/^\d{4}/)) return new Date(dateStr);

      const [month, day] = dateStr.split('-').map(Number);
      return new Date(currentYear, month - 1, day);
    },
  },
};
