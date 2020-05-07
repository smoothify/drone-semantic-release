const { DRY_RUN, GIT_ASSETS, GIT_COMMIT_TEMPLATE, GENERATE_CHANGELOG, GENERATE_PACKAGEJSON, VERSION_FILENAME} = process.env;

let plugins = [
  '@semantic-release/commit-analyzer',
  '@semantic-release/release-notes-generator'
];

if (GENERATE_CHANGELOG === "true") {
  plugins.push('@semantic-release/changelog');
}

if (GENERATE_PACKAGEJSON === "true") {
  plugins = [
    ...plugins,
    ["@semantic-release/npm", {
      "npmPublish": false
    }],
  ]
}

let execCommandName = DRY_RUN === "true" ? "verifyReleaseCmd" : "prepareCmd";
let execCommand = {};
execCommand[execCommandName] = "echo \"${nextRelease.version}\" > " + VERSION_FILENAME || ".release-version";

module.exports = {
  branches: [
    '+([0-9])?(.{+([0-9]),x}).x',
    'master',
    'next',
    'next-major',
    {name: 'beta', prerelease: true},
    {name: 'alpha', prerelease: true},
    {name: 'dev', prerelease: true}
  ],
  plugins: [
    ...plugins,
    ["@semantic-release/exec", execCommand],
    ["@semantic-release/git", {
      "assets": GIT_ASSETS ? GIT_ASSETS.split(",") : ["CHANGELOG.md", "package.json"],
      "message": GIT_COMMIT_TEMPLATE || "chore(release): ${nextRelease.version} [skip ci]"
    }]
  ],
};
