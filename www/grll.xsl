<?xml version="1.0" encoding="utf-8" ?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" >
<xsl:output omit-xml-declaration="yes" indent="yes"/>

<xsl:template match="/GrllTestRun">

    <xsl:variable name="link">
        <xsl:value-of select="link" />
    </xsl:variable>

    <xsl:variable name="verdict">
        <xsl:choose>
        <xsl:when test="count(results/status[text() = 'failed']) = 0">passed</xsl:when>
        <xsl:otherwise>failed</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>


<html><head><title>::Grll::</title>
<link href="/css/grll.css" rel='stylesheet' type='text/css'/>
</head>
<body>
    <h1 class="{$verdict}"><xsl:value-of select="run"/></h1>

    <h4><span><xsl:value-of select="ts"/></span>
        <span><a href="{$link}">artifacts</a></span>
    </h4>

    <h4>
        <span class="passed">Passed = <xsl:value-of select="count(results/status[text() = 'passed'])"/></span>
        <span class="failed">Failed = <xsl:value-of select="count(results/status[text() = 'failed'])"/></span>
        <span class="skipped">Skipped = <xsl:value-of select="count(results/status[text() = 'skipped'])"/></span>

    </h4>

    <div class="object">
        <h2 class="h2">Results</h2>
        <table cellspacing="0">
        <tr class="header"><td>Name</td><td>Status</td><td>Message</td></tr>
        <xsl:apply-templates select='results' />
        </table>
    </div>

    <div class="object">
        <h2 class="h2">Measurements</h2>
        <table cellspacing="0">
            <tr class="header">
                <td>Name</td>
                <td>Value</td>
                <td>Unit</td>
            </tr>
            <xsl:apply-templates select='values' />
        </table>
    </div>

    <div class="object">
        <h2 class="h2">tags</h2>
        <div class="taglist">
        <ul id="tags">
            <xsl:apply-templates select='tags' />
        </ul>
        </div>
    </div>

</body>
</html>
</xsl:template>


<xsl:template match='results'>
    <xsl:variable name="class">
        <xsl:choose>
        <xsl:when test="normalize-space(./status) = 'passed'">passed</xsl:when>
        <xsl:when test="normalize-space(./status) = 'failed'">failed</xsl:when>
        <xsl:when test="normalize-space(./status) = 'skipped'">skipped</xsl:when>
        <xsl:otherwise>failed</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>
    <xsl:variable name="oddeven">
        <xsl:choose>
        <xsl:when test="position() mod 2 = 0">even</xsl:when>
        <xsl:otherwise>odd</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>
    <xsl:variable name="testname"><xsl:value-of select="normalize-space(./test)"/></xsl:variable>
    <tr class="{$oddeven}">
        <td class="{$class}">
            <a href="/api/testrun?test={$testname}">
            <xsl:copy-of select="$testname" />
            </a>
        </td>
        <td class="{$class}"><xsl:value-of select="normalize-space(./status)"/></td>
        <td class="{$class}"><xsl:value-of select="normalize-space(./msg)"/></td>
    </tr>
</xsl:template>

<xsl:template match='values'>
    <xsl:variable name="oddeven">
        <xsl:choose>
        <xsl:when test="position() mod 2 = 0">even</xsl:when>
        <xsl:otherwise>odd</xsl:otherwise>
        </xsl:choose>
    </xsl:variable>
    <tr class="{$oddeven}">
        <td><xsl:value-of select="normalize-space(./test)"/></td>
        <td><xsl:value-of select="normalize-space(./value)"/></td>
        <td><xsl:value-of select="normalize-space(./unit)"/></td>
    </tr>
</xsl:template>


<xsl:template match='tags'>
        <li><a href="#">
            <xsl:value-of select="normalize-space(.)"/>
        </a></li>
</xsl:template>


</xsl:stylesheet>

