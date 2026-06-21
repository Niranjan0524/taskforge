const Footer=()=>{


    return <>
    <footer className="mt-20 border-t border-black/10 relative overflow-hidden">
  {/* Glow */}
  <div className="absolute right-20 bottom-0 h-56 w-56 rounded-full bg-orange-500/20 blur-3xl" />

  <div className="max-w-7xl mx-auto px-6 py-16">
    <div className="grid grid-cols-1 md:grid-cols-4 gap-12">

      {/* Brand */}
      <div className="md:col-span-2">
        <div className="flex items-center gap-3 mb-4">
          <div className="h-11 w-11 rounded-xl bg-orange-500/15 border border-orange-500/30 flex items-center justify-center">
            ⚙️
          </div>

          <h2 className="text-3xl font-bold text-[#191a16]">
            TaskForge
          </h2>
        </div>

        <p className="max-w-md text-[#1c1d19]/70 leading-relaxed">
          Reliable background job processing, worker orchestration,
          queue monitoring and real-time task tracking built for scale.
        </p>
      </div>

      {/* Navigate */}
      <div>
        <h4 className="text-xs tracking-[0.25em] font-semibold text-[#1c1d19]/50 mb-5">
          NAVIGATE
        </h4>

        <div className="space-y-3">
          <a href="/" className="block hover:text-orange-600 transition">
            Dashboard
          </a>
          <a href="/tasks" className="block hover:text-orange-600 transition">
            Tasks
          </a>
          <a href="/workers" className="block hover:text-orange-600 transition">
            Workers
          </a>
          <a href="/analytics" className="block hover:text-orange-600 transition">
            Analytics
          </a>
        </div>
      </div>

      {/* Connect */}
      <div>
        <h4 className="text-xs tracking-[0.25em] font-semibold text-[#1c1d19]/50 mb-5">
          CONNECT
        </h4>

        <div className="flex flex-wrap gap-3">
          <a
            href="https://github.com/Niranjan0524"
            className="px-4 py-2 rounded-full backdrop-blur-md bg-white/40 border border-black/10 hover:border-orange-500/40 hover:bg-orange-500/10 transition"
          >
            GitHub
          </a>

          <a
            href="www.linkedin.com/in/niranjan05"
            className="px-4 py-2 rounded-full backdrop-blur-md bg-white/40 border border-black/10 hover:border-orange-500/40 hover:bg-orange-500/10 transition"
          >
            LinkedIn
          </a>

          <a
            href="https://niranjan5.me"
            className="px-4 py-2 rounded-full backdrop-blur-md bg-white/40 border border-black/10 hover:border-orange-500/40 hover:bg-orange-500/10 transition"
          >
            Portfolio
          </a>
        </div>
      </div>
    </div>

    {/* Bottom */}
    <div className="mt-12 pt-6 border-t border-black/10 flex flex-col md:flex-row justify-between gap-4 text-sm text-[#1c1d19]/55">
      <span>© 2026 TaskForge — All rights reserved.</span>

      <span>
        Built with
        <span className="text-orange-700 font-semibold"> Go </span>
        ,
        <span className="text-orange-700 font-semibold"> React</span>
        &
        <span className="text-orange-700 font-semibold"> Redis </span>
      </span>
    </div>
  </div>
</footer>
        </>
}

export default Footer;