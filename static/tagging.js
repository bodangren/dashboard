'use strict';

const TagManager = (function() {
  const STORAGE_KEY = 'dashboard-tags';
  const VERSION_KEY = 'dashboard-tags-version';
  const CURRENT_VERSION = 1;

  function getRawTags() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (!raw) return {};
      return JSON.parse(raw);
    } catch (e) {
      return {};
    }
  }

  function saveTags(tags) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(tags));
    localStorage.setItem(VERSION_KEY, String(CURRENT_VERSION));
  }

  function getTagsForRepo(repoPath) {
    const tags = getRawTags();
    return tags[repoPath] || [];
  }

  function setTagsForRepo(repoPath, tagList) {
    const tags = getRawTags();
    if (!tagList || tagList.length === 0) {
      delete tags[repoPath];
    } else {
      tags[repoPath] = tagList.slice();
    }
    saveTags(tags);
  }

  function addTag(repoPath, tag) {
    const current = getTagsForRepo(repoPath);
    const trimmed = tag.trim();
    if (!trimmed || current.includes(trimmed)) return current;
    const updated = current.concat([trimmed]);
    setTagsForRepo(repoPath, updated);
    return updated;
  }

  function removeTag(repoPath, tag) {
    const current = getTagsForRepo(repoPath);
    const trimmed = tag.trim();
    const updated = current.filter(function(t) { return t !== trimmed; });
    setTagsForRepo(repoPath, updated);
    return updated;
  }

  function getAllTags() {
    const tags = getRawTags();
    const all = [];
    for (const repoPath in tags) {
      for (const tag of tags[repoPath]) {
        if (all.indexOf(tag) === -1) {
          all.push(tag);
        }
      }
    }
    return all.sort();
  }

  function getReposWithTag(tag) {
    const tags = getRawTags();
    const matches = [];
    const trimmed = tag.trim();
    for (const repoPath in tags) {
      if (tags[repoPath].indexOf(trimmed) !== -1) {
        matches.push(repoPath);
      }
    }
    return matches;
  }

  function clearAllTags() {
    localStorage.removeItem(STORAGE_KEY);
    localStorage.removeItem(VERSION_KEY);
  }

  return {
    getTagsForRepo: getTagsForRepo,
    setTagsForRepo: setTagsForRepo,
    addTag: addTag,
    removeTag: removeTag,
    getAllTags: getAllTags,
    getReposWithTag: getReposWithTag,
    clearAllTags: clearAllTags,
    STORAGE_KEY: STORAGE_KEY
  };
})();